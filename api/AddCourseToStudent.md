**Add Course to Student**
----
  Adds the course with given id to the student with the given username.

* **URL**

  /students/:student_username/courses/:course_id

* **Method:**

  `PUT`
  
*  **URL Params**

   **Required:**
 
   `student_username=[string]`<br/>
   `course_id=[string]`
   
* **Data Params**

    None

* **Success Response:**

  * **Code:** 200 OK <br />
    **Content:** `{
    "courses": [
        "5cdc6d2339f578a97998a8d9",
        "5cdc6d7439f578a97998a8da",
        "5cdc6dac39f578a97998a8db",
        "5cdc6dfe39f578a97998a8dc"
    ],
    "id": "5cdfd0da04e5d52398f3d8f9",
    "name": "John Doe",
    "username": "student_username"
}` (This is the updated student with the identifiers of the attended courses)
 
* **Error Response:**

  * **Code:** 409 CONFLICT <br />
    **Content:** `{ error : "Conflict - The resource already exists"}`
    This is returned when the client requires to append a course that the user has already subscribed to
    
  OR

  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Internal Server Error" }`
       
  OR

  * **Code:** 404 NOT FOUND <br />
    **Content:** `{ error : "Course Not Found" }`
    This is returned where a course with the given id does not exist
    
  OR

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{ error : "No token provided" }`
    
  OR

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{ error : "Wrong credentials" }`
    
  OR

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{ error : "Permission denied" }` This error may occur as the course creation is allowed to teachers only.
    
  OR

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{ error : "Expired token" }`
     