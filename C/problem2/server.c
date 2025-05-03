/* https://protohackers.com/problem/2
 * cc -oserver server.c -l sqlite3 && ./server 8007
 */

#include <byteswap.h>
#include <netinet/in.h>
#include <sqlite3.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

static sqlite3* create_db(void)
{
    sqlite3 *db;
    char *db_err;
    int db_ret;

    if (sqlite3_open(":memory:", &db) != SQLITE_OK) {
        fprintf(stderr, "sqlite3_open error: %s\n", sqlite3_errmsg(db));
        exit(EXIT_FAILURE);
    }

    db_ret = sqlite3_exec(
        db, "create table prices(timestamp integer not null unique primary key, price integer not null);",
        NULL, NULL, &db_err);
    if (db_ret != SQLITE_OK) {
        fprintf(stderr, "sqlite3_exec create table: %s\n", db_err);
        sqlite3_free(db_err);
        exit(EXIT_FAILURE);
    }

    return db;
}

static void destroy_db(sqlite3 *db)
{
    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "sqlite3_close error: %s\n", sqlite3_errmsg(db));
        exit(EXIT_FAILURE);
    }
}

static void insert_price(sqlite3 *db, int timestamp, int price)
{
    char *db_sql, *db_err;

    db_sql = sqlite3_mprintf(
        "insert or ignore into prices(timestamp, price) values(%d, %d);",
        timestamp, price);

    if (sqlite3_exec(db, db_sql, NULL, NULL, &db_err) != SQLITE_OK) {
        fprintf(stderr, "sqlite3_exec insert: %s\n", db_err);
        sqlite3_free(db_err);
        exit(EXIT_FAILURE);
    }

    sqlite3_free(db_sql);
}

static int query_price(sqlite3 *db, int timestamp_from, int timestamp_to)
{
    char *db_sql;
    int db_ret, mean_price;
    sqlite3_stmt *db_stmt;

    db_sql = sqlite3_mprintf( /* indeed, I could have used sqlite3_bind_int instead */
        "select round(coalesce(avg(price),0)) from prices where timestamp between %d and %d;",
        timestamp_from, timestamp_to);

    if ((db_ret = sqlite3_prepare_v2(db, db_sql, -1, &db_stmt, NULL)) != SQLITE_OK) {
        fprintf(stderr, "sqlite3_prepare_v2: %s\n", sqlite3_errstr(db_ret));
        exit(EXIT_FAILURE);
    }

    if ((db_ret = sqlite3_step(db_stmt)) != SQLITE_ROW) {
        fprintf(stderr, "sqlite3_step: %s\n", sqlite3_errstr(db_ret));
        exit(EXIT_FAILURE);
    }

    mean_price = sqlite3_column_int(db_stmt, 0);

    sqlite3_finalize(db_stmt);
    sqlite3_free(db_sql);

    return mean_price;
}

static int read_data(int conn, char *data, int max_data_len)
{
    int read_len, data_len = 0;

    do {
        if ((read_len = read(conn, data + data_len, max_data_len - data_len)) < 0) {
            perror("read error");
            exit(EXIT_FAILURE);
        }

        if (read_len == 0) {
            break;
        }

        data_len += read_len;
    }
    while (data_len < max_data_len);

    return data_len;
}

static void handle_connection(int conn)
{
    sqlite3 *db;
    union data
    {
        char raw[9];
        struct msg
        {
            char type;
            unsigned int int1;
            unsigned int int2;
        } __attribute__((packed)) msg;
    } __attribute__((packed)) data;
    int data_len;

    db = create_db();

    while ((data_len = read_data(conn, data.raw, sizeof(data))) == sizeof(data)) {
        data.msg.int1 = __bswap_32(data.msg.int1);
        data.msg.int2 = __bswap_32(data.msg.int2);

        switch (data.msg.type) {
            case 'I':
                insert_price(db, data.msg.int1, data.msg.int2);
                break;

            case 'Q':
                int mean_price = query_price(db, data.msg.int1, data.msg.int2);
                mean_price = __bswap_32(mean_price);

                if (send(conn, &mean_price, sizeof(mean_price), 0) < 0) {
                    perror("send error");
                    exit(EXIT_FAILURE);
                }
                break;
        }
    }

    destroy_db(db);

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
