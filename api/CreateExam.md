**Create Course**
----
  Adds an exam with the information provided in the JSON body of the request.
  A JWT token has to be provided to authenticate the request and verify that
  the user has the authorization to perform the operation.

* **URL**

  /exams/

* **Method:**

  `POST`
  
*  **URL Params**

   **Required:**
 
   None
   

* **Data Params**

    `{course: "IdCourse", call: "2", date: "21-03-2019", startTime: "10:30", room: "A1",
    	  expirationDate: "18-03-2019"}`

* **Success Response:**

  * **Code:** 201 CREATED <br />
    **Content:** `{ id: "5ce0165fe2c5c2136899fad5", course: "IdCourse", 
                                     call: "2", date: "21-03-2019", startTime: "10:30",
                                     room: "A1", expirationDate: "18-03-2019", students: []}`
 
* **Error Response:**

  * **Code:** 409 CONFLICT <br />
    **Content:** `{ error : "Conflict - The resource already exists"}`
    This is returned when an exam for the given course with the same call already exists

  OR

  * **Code:** 400 BAD REQUEST <br />
    **Content:** `{ error : "Bad request" }`
    
  OR

  * **Code:** 400 BAD REQUEST <br />
    **Content:** `{ error : "Malformed token" }`
    
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
    **Content:** `{ error : "Permission denied" }` This error may occur as the course creation is allowed to teachers only.
    
  OR

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{ error : "Expired token" }`