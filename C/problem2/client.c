/* https://protohackers.com/problem/2
 * cc -oclient client.c && ./client 8007
 */

#include <arpa/inet.h>
#include <byteswap.h>
#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

static void send_data(int sock, char *data, int data_len)
{
    if (send(sock, data, data_len, 0) < 0) {
        perror("send error");
        exit(EXIT_FAILURE);
    }
}

int main(int argc, char *argv[])
{
    int sock;
    struct linger sock_opt_linger;
    struct sockaddr_in addr;

    char data[9];
    int data_len;
    int res;

    if (argc != 2) {
        printf("missing or wrong argument(s): client <port number>\n");
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

    memcpy(data, "\x49\x00\x00\x30\x39\x00\x00\x00\x65", 9);
    send_data(sock, data, sizeof(data));

    memcpy(data, "\x49\x00\x00\x30\x3a\x00\x00\x00\x66", 9);
    send_data(sock, data, sizeof(data));

    memcpy(data, "\x49\x00\x00\x30\x3b\x00\x00\x00\x64", 9);
    send_data(sock, data, sizeof(data));

    memcpy(data, "\x49\x00\x00\xa0\x00\x00\x00\x00\x05", 9);
    send_data(sock, data, sizeof(data));

    /* same timestamp */
    memcpy(data, "\x49\x00\x00\xa0\x00\x00\x00\x00\x06", 9);
    send_data(sock, data, sizeof(data));

    /* invalid type */
    memcpy(data, "\x58\x00\x00\x00\x00\x00\x00\x00\x00", 9);
    send_data(sock, data, sizeof(data));

    memcpy(data, "\x51\x00\x00\x30\x00\x00\x00\x40\x00", 9);
    send_data(sock, data, sizeof(data));

    if ((read(sock, &res, sizeof(res))) < 0) {
        perror("read error");
        exit(EXIT_FAILURE);
    }

    printf("mean price: %d - 0x%08x\n", __bswap_32(res), __bswap_32(res));
    fflush(stdout);

    if (shutdown(sock, SHUT_WR) < 0) {
        perror("shutdown error");
        exit(EXIT_FAILURE);
    }

    if (close(sock) < 0) {
        perror("close error");
        exit(EXIT_FAILURE);
    }

    return EXIT_SUCCESS;
}
