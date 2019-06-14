**Publish Course Notification**
----
  Publish a notification about the course with the given id on 
  a message queue that will be consumed by a notification management
  microservice, which is responsible for sending the notification via e-mail

* **URL**

  /notification/course/:course_id

* **Method:**

  `POST`
  
*  **URL Params**

   **Required:**
 
   `course_id=[string]`
   

* **Data Params**

    `{"message":"The lecture has been moved..."}`

* **Success Response:**

  * **Code:** 200 OK <br />
 
* **Error Response:**

  * **Code:** 404 NOT FOUND <br />
    **Content:** `{ error : "Course Not Found"}`

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
    **Content:** `{ error : "Permission denied" }`
    
  OR

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{ error : "Expired token" }`