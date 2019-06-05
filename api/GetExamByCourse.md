**Get Course By**
----
    Finds exams by id of the course they belong to.
* **URL**

  /exams/:course
  
* **Method:**

  `GET`
  
*  **URL Params**

   **Required:**
 
   `course=[string]`
   
   `course` is the id of the course the exams to find belong to 


* **Data Params**

    None

* **Success Response:**

  * **Code:** 200 OK <br />
    **Content:** `[
    {   "call": 1,
        "course": "5cdff2e80cdbd6af9de9abf3",
        "date": "1-9-2019",
        "expirationDate": "15-8-2019",
        "id": "5ce17880c0fd26b2d3a6e233",
        "room": "A1",
        "startTime": "10:00",
        "students": []},
    {   "call": 2,
        "course": "5cdff2e80cdbd6af9de9abf3",
        "date": "15-9-2019",
        "expirationDate": "1-9-2019",
        "id": "5ce17b39c0fd26b2d3a6e234",
        "room": "B2",
        "startTime": "10:00",
        "students": []}
    ]`
 
* **Error Response:**

  * **Code:** 404 NOT FOUND <br />
    **Content:** `{ error : "Not found"}`

  OR

  * **Code:** 400 BAD REQUEST <br />
    **Content:** `{ error : "Bad request" }`
    
    
  OR

  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Internal Server Error" }`
    
  OR

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{ error : "No token provided" }`
    
  OR

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{ error : "Wrong credentials" }`
    
  OR

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{ error : "Expired token" }`