POST Register

Request
```
{
    "name": "Ivan",
    "surname": "Petrov",
    "gender": "male",
    "birth_date": "21-03-1992",
    "height_cm": 0,
    "weight_kg": 0,
    "sport_activity_level_id": 1,
    "town_id": 1,
    "phone_number": "+51",
    "email": "1113e@example.com",
    "password": "strongPassword123",
    
    "is_have_injury": false,
    "injury_description": "lol", // nullable
    "photo": "https://example.com/photos/ivan.jpg" // nullable
}
```

Response
````
{
    "access_token": "...." ,
    "refresh_token": "...",
}
````
---

