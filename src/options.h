/**
 * @package RER.
 */

#include <stdbool.h>

typedef struct {
    char date[11];
    char currency[4];
    float amount;
    bool inverse;
    bool raw;
} cli_options;

const cli_options parse_options(int argc, char **argv);
