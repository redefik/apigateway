**Create Course**
----
  Adds a course with the information provided in the JSON body of the request.
  A JWT token has to be provided to authenticate the request and verify that
  the user has the authorization to perform the operation.

* **URL**

  /courses/

* **Method:**

  `POST`
  
*  **URL Params**

   **Required:**
 
   None
   

* **Data Params**

    `{name:"Advanced Calculus", department:"Science", teacher: "Doe", "year":"2019-2020", semester:2,
	  description:"Limits, Derivatives, Integrals",
	  schedule:[{day:"lun", startTime: "10:00", endTime: "11:00", room: "A4" },
				{day: "mar", startTime: "11:30", endTime: "12:30", room: "B9"}]
    }`

* **Success Response:**

  * **Code:** 201 CREATED <br />
    **Content:** `{id:"5cda791f5aec95bb5a5abd7c",
                   name:"Advanced Calculus", department:"Science", teacher: "Doe", year: "2019-2020", semester: 2,
	               description: "Limits, Derivatives, Integrals",
	               schedule:[{day: "lun", startTime: "10:00", endTime: "11:00", room: "A4" },
				                {day: "mar", startTime: "11:30", endTime: "12:30", room: "B9"}]
                     }`
 
* **Error Response:**

  * **Code:** 409 CONFLICT <br />
    **Content:** `{ error : "Conflict - The resource already exists"}`
    This is returned when a course with the given name already exists

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