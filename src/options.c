/**
 * @package RER.
 */

#include <stdlib.h>
#include <stdio.h>
#include <getopt.h>
#include <stdbool.h>
#include <string.h>
#include <ctype.h>
#include "options.h"

#define DATE_LENGTH 10 // strlen(dd.mm.yyyyy)
#define VERSION "rer version 1.0.0"

void usage() {
    puts("usage: rer [option]... \n");
    puts("a command line tool to obtain ruble exchnge rate.");
    puts("  -h, --help      display this help and exit.");
    puts("  -v, --version   display the version number.");
    puts("  -d, --date      a date in format dd.mm.yyyy for wich to display exchange rate (default: current date).");
    puts("  -c, --currency  a currency for wich to display rate (default: all).");
    puts("  -r, --raw       print rates without formatting.");
    puts("  -a, --amount    an amount of currency calculate the rate (default: 1). implicitly sets --raw option.");
    puts("  -i, --inverse   display inverse currency rate. implicitly sets --raw option.");
    puts("");
}

void parse_date(const char* input, char* date)
{
    if (strlen(input) != DATE_LENGTH) {
        fprintf(stderr, "wrong date '%s'.\n", input);
        exit(EXIT_FAILURE);
    }

    unsigned short int day, month, year;
    if (sscanf(input, "%2hu.%2hu.%4hu", &day, &month, &year) != 3) {
        fprintf(stderr, "wrong date '%s'.\n", input);
        exit(EXIT_FAILURE);
    }
    if (day > 31  || month > 12) {
        fprintf(stderr, "wrong date '%s'.\n", input);
        exit(EXIT_FAILURE);
    }

    sprintf(date, "%02d/%02d/%d", day, month, year);
}

const cli_options parse_options(int argc, char **argv)
{
    cli_options options = {
        .date = "---.--.--",
        .currency = "---",
        .amount =  1.0,
        .inverse = false,
        .raw = false,
    };

    static struct option long_options[] =
    {
        {"date",	   required_argument, NULL, 'd'},
        {"currency", required_argument, NULL, 'c'},
        {"amount",	 required_argument, NULL, 'a'},
        {"inverse",	 no_argument,       NULL, 'i'},
        {"raw",	     no_argument,       NULL, 'r'},
        {"help",	   no_argument,       NULL, 'h'},
        {"version",	 no_argument,       NULL, 'v'},
        {NULL, 0, NULL, 0}
    };

    int ch;
    while ((ch = getopt_long(argc, argv, ":a:c:d:irhv", long_options, NULL)) != -1) {
        switch (ch) {
            case 'd':
                parse_date(optarg, options.date);
                break;

            case 'c':
                if (strlen(optarg) > 3) {
                    fprintf(stderr, "wrong currency code '%s'.\n", optarg);
                    exit(EXIT_FAILURE);
                }
                options.currency[0] = toupper(optarg[0]);
                options.currency[1] = toupper(optarg[1]);
                options.currency[2] = toupper(optarg[2]);
                break;

            case 'a':
                options.amount = atof(optarg);
                options.raw = 1;
                break;

            case 'i':
                options.inverse = 1;
                options.raw = 1;
                break;

            case 'r':
                options.raw = 1;
                break;

            case 'h':
                usage();
                exit(EXIT_SUCCESS);
                break;

            case 'v':
                puts(VERSION);
                exit(EXIT_SUCCESS);
                break;

            case ':':
                fprintf(stderr, "'%c' requires an argument\n", optopt);
                exit(EXIT_FAILURE);
                break;

            case '?':
            default :
                fprintf(stderr, "unknown option '%c'.\n", optopt);
                exit(EXIT_FAILURE);
        }
    }

    return options;
}
