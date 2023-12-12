# Test - Cosmart Software Engineers
Test software engineer - cosmart

#### How running this code
    please before run using makefile see makefile for details.
    command:
    make run/service => (start service)
    make clean/cache => (clean cache)
    make run/download => (download package)
    make clean/package => (remove package)
    make test/coverage => (coverage unit testing)

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
