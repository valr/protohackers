/* https://protohackers.com/problem/1
 * cc -oserver server.c -lm `pkg-config --cflags --libs jansson` && ./server 8007
 */

#include <jansson.h>
#include <math.h>
#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

#define MAXDATALEN 32768

static int number_is_prime(double number)
{
    if (number != round(number) || number < 2.0) {
        return 0;
    }

    for (double i = 2.0; i * i <= number; i = i + 1.0) {
        if (fmod(number, i) == 0.0) {
            return 0;
        }
    }

    return 1;
}

static int process_request(const char *request, char *reply, size_t reply_size)
{
    json_auto_t *json_request, *json_reply;
    json_t *json_method = NULL, *json_number = NULL;
    json_error_t error;

    if ((json_request = json_loads(request, JSON_DECODE_INT_AS_REAL, &error)) != NULL) {
        json_method = json_object_get(json_request, "method");
        json_number = json_object_get(json_request, "number");
    }

    if (json_is_string(json_method) != 0 && strcmp(json_string_value(json_method), "isPrime") == 0 &&
        json_is_real(json_number) != 0) {
        /* conforming request:
         * has field "method" containing "isPrime"
         * has field "number" containing valid JSON number */

        /* conforming reply:
         * has field "method" containing "isPrime"
         * has field "prime" containing boolean (true or false) */

        int prime = number_is_prime(json_real_value(json_number));

        json_reply = json_pack("{s, s, s, b}", "method", "isPrime", "prime", prime);
        strcpy(&reply[json_dumpb(json_reply, reply, reply_size, 0)],"\n\0");

        printf("conforming request: %s - %s\n", request, reply);
        fflush(stdout);

        return 0;
    } else {
        /* malformed request:
         * invalid JSON object
         * missing field method/number
         * invalid method/number */

        /* malformed reply (simplified):
         * {"method": "isPrime", "error": "invalid request"} */

        json_reply = json_pack("{s, s, s, s}", "method", "isPrime", "error", "invalid request");
        strcpy(&reply[json_dumpb(json_reply, reply, reply_size, 0)],"\n\0");

        printf("malformed request: %s - %s\n", request, reply);
        fflush(stdout);

        return -1;
    }

    /* return conforming or malformed reply
     * which is a single line containing JSON object terminated by newline */
}

static void handle_connection(int conn)
{
    FILE *file;
    char request[MAXDATALEN], reply[MAXDATALEN];
    int request_status;

    if ((file = fdopen(dup(conn), "r")) == NULL) {
        perror("fdopen error");
        exit(EXIT_FAILURE);
    }

    while(fgets(request, sizeof(request), file) != NULL) {
        request_status = process_request(request, reply, sizeof(reply));

        if (send(conn, reply, strlen(reply), 0) < 0) {
            perror("send error");
            exit(EXIT_FAILURE);
        }

        if (request_status == -1) {
            break;
        }
    }

    if (fclose(file) != 0) {
        perror("fclose error");
        exit(EXIT_FAILURE);
    }

    if (shutdown(conn, SHUT_RDWR) < 0) {
        perror("shutdown error");
        exit(EXIT_FAILURE);
    }

    if (close(conn) < 0) {
        perror("close error");
        exit(EXIT_FAILURE);
    }
}

int main(int argc, char *argv[])
{
    int sock, conn;
    struct linger sock_opt_linger;
    struct sockaddr_in addr;
    int addr_len = sizeof(addr);
    int port;
    pid_t pid;

    if (argc != 2) {
        printf("missing or wrong argument(s): server <port number>\n");
        exit(EXIT_FAILURE);
    }

    if ((sock = socket(AF_INET, SOCK_STREAM, 0)) < 0) {
        perror("socket error");
        exit(EXIT_FAILURE);
    }

    sock_opt_linger.l_onoff = 1;
    sock_opt_linger.l_linger = 30;

    if ((setsockopt(sock, SOL_SOCKET, SO_LINGER, &sock_opt_linger, sizeof(sock_opt_linger))) < 0) {
        perror("setsockopt error");
        exit(EXIT_FAILURE);
    }

    bzero(&addr, sizeof(addr));
    addr.sin_family = AF_INET;
    addr.sin_addr.s_addr = INADDR_ANY;
    addr.sin_port = htons(atoi(argv[1]));

    if (bind(sock, (struct sockaddr*)&addr, sizeof(addr)) < 0) {
        perror("bind error");
        exit(EXIT_FAILURE);
    }

    if (listen(sock, 5) < 0) {
        perror("listen error");
        exit(EXIT_FAILURE);
    }

    while ((conn = accept(sock, (struct sockaddr*)&addr, (socklen_t*)&addr_len)) >= 0) {
        pid = fork();
        switch(pid) {
            case -1:
                perror("fork error");
                exit(EXIT_FAILURE);
                break;
            case 0:
                handle_connection(conn);
                exit(EXIT_SUCCESS);
                break;
        }
    }

    if (conn < 0) {
        perror("accept error");
        exit(EXIT_FAILURE);
    }

    if (close(sock) < 0) {
        perror("close error");
        exit(EXIT_FAILURE);
    }

    return EXIT_SUCCESS;
}
