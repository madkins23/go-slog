# `server`

The `cmd/server` application displays benchmark and verification results as web pages.

## Server Pages

The home page shows links to various test data pages and the warnings:
![The home page shows links to various test data pages and the warnings.](images/home.png)

Test pages show the same tables as `tabular` plus charts comparing the results:
![Test pages show the same tables as `tabular` plus charts comparing the results.](images/test.png)

Handler pages show similar tables plus charts comparing the results:
![Handler pages show similar tables plus charts comparing the results.](images/handler.png)

The scores page shows how different handlers related on a functionality vs. performance chart:
![Scores page shows how different handlers related on a functionality vs. performance chart](images/scores.png)

The warning page shows all the defined warnings with descriptions:
![The warning page shows all the defined warnings with descriptions.](images/warnings.png)

## GitHub Pages

Once a week (or whenever code is committed to the `go-slog` repository)
the server is run and all pages are copied to the `docs` directory.
This is committed back into the repository and
[GitHub Pages](https://pages.github.com/) serves the
[recent benchmark data](https://madkins23.github.io/go-slog/index.html).
