/* https://protohackers.com/problem/0
 * cc -oserver server.c && ./server 8007 */

#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

static void handle_connection(int conn)
{
    char data[1024];
    int data_len;

    while ((data_len = read(conn, data, sizeof(data))) > 0) {
        if (send(conn, data, data_len, 0) < 0) {
            perror("send error");
            exit(EXIT_FAILURE);
        }
    }

    if (data_len < 0) {
        perror("read error");
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
