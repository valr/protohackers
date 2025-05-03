/* https://protohackers.com/problem/1
cc -oclient client.c && ./client 8007 '{"method":"isPrime","number":3}
{"method":"isPrime","number":643576916629744492749636813248163978082647287841524867798}
{"method":"isPrime","number":4}
{"method":"isPrime","number":1.4}
{"method":"isPrime","number":false}
{"method":"isPrime","number":123}'
*/

#include <arpa/inet.h>
#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

int main(int argc, char *argv[])
{
    int sock;
    struct linger sock_opt_linger;
    struct sockaddr_in addr;

    char data[1024];
    int data_len;

    if (argc != 3) {
        printf("missing or wrong argument(s): client <port number> <data to send>\n");
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
    addr.sin_port = htons(atoi(argv[1]));

    if (inet_pton(AF_INET, "127.0.0.1", &addr.sin_addr) <= 0) {
        perror("inet_pton error");
        exit(EXIT_FAILURE);
    }

    if ((connect(sock, (struct sockaddr*)&addr, sizeof(addr))) < 0) {
        perror("connect error");
        exit(EXIT_FAILURE);
    }

    if (send(sock, argv[2], strlen(argv[2]) + 1, 0) < 0) {
        perror("send error");
        exit(EXIT_FAILURE);
    }

    if (shutdown(sock, SHUT_WR) < 0) {
        perror("shutdown error");
        exit(EXIT_FAILURE);
    }

    while ((data_len = read(sock, data, sizeof(data))) > 0) {
        fwrite(data, data_len, 1, stdout);
        fflush(stdout);
    }

    if (data_len < 0) {
        perror("read error");
        exit(EXIT_FAILURE);
    }

    if (close(sock) < 0) {
        perror("close error");
        exit(EXIT_FAILURE);
    }

    return EXIT_SUCCESS;
}
