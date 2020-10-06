/**
 * @package RER.
 */

#include <stdio.h>
#include <stdbool.h>
#include "options.h"
#include "fetch.h"
#include "decoder.h"

#define BASE_URL "https://www.cbr.ru/scripts/XML_daily.asp"

int main(int argc, char **argv)
{
    const cli_options options = parse_options(argc, argv);

    char url[65];
    if (options.date[0] != '-') {
        sprintf(url, "%s?date_req=%s", BASE_URL, options.date);
    }
    else {
        strcpy(url, BASE_URL);
    }

    char *encoded_xml = NULL;
    encoded_xml = fetch(url);

    rate rate;
    parse_rates(encoded_xml, &rate);

    if (!options.raw) {
        printf("Exchange rate of the ruble, %s\n", rate.date);
        puts("--------------------------------------");
    }

    unsigned short i;
    float value;
    for (i = 0; i < rate.total; i++) {
        if (options.currency[0] != '-' && strcmp(options.currency, rate.records[i].code) != 0) {
            continue;
        }
        value = rate.records[i].value;
        if (options.raw) {
            if (options.inverse) {
                value = 1 / value;
            }
            value = options.amount * value / rate.records[i].nominal;
            printf("%.4f\n", value);
        }
        else {
            printf("%9.4f руб. = %d %s (%s)\n", value, rate.records[i].nominal, rate.records[i].name, rate.records[i].code);
        }
    }

    return EXIT_SUCCESS;
}
