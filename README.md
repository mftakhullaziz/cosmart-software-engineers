# Test - Cosmart Software Engineers
Test software engineer - cosmart

#### How running this code
    if want running using makefile this docs for use, please before run see makefile for details.
    start service: make run/service
    clean cache: make clean/cache
    download package: make run/download
    remove package: make clean/package
    coveraget unit testing: make test/coverage

#### API Curl
    Get Books By Genre:
    curl --location 'http://localhost:8080/books/love'

    sample response:
    {
    "status": "200 OK",
    "is_success": true,
    "message": "fetch data books successfully!",
    "total_data": 2,
    "data": [
        {
            "title": "Wuthering Heights",
            "author": [
                "Emily Brontë"
            ],
            "edition_number": 2123
        },
        {
            "title": "Rose in Bloom",
            "author": [
                "Louisa May Alcott",
                "Harriet Roosevelt Richards"
            ],
            "edition_number": 257
        }
    ]

    Save Books Pick Up Schedule
    curl --location 'http://localhost:8080/books/schedule' \
    --header 'Content-Type: application/json' \
    --data '{
        "book_info": {
            "title": "C programming phase 1",
            "author": [
                "author2",
                "author3"
            ],
            "edition_number": 1
        },
        "pick_up_date": "2023-12-01",
        "genre": "love"
    }'

    sample response:
    {
    "status": "201 CREATED",
    "is_success": true,
    "message": "save new data books successfully!",
    "total_data": 1,
    "data": {
        "book_info": {
            "title": "C programming phase 1",
            "author": [
                "author2",
                "author3"
            ],
            "edition_number": 1
        },
        "pick_up_date": "2023-12-01",
        "genre": "love"
        }
    }
