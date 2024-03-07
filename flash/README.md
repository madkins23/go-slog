# Flash Handler

## Flash over Sloggy

* Replace `bytes.Buffer` w/`addXXX()`
* Reuse pools for:
    - Log record output
    - Basic attribute list
    - Source data record
    - Composers
