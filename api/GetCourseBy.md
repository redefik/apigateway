**Get Course By**
----
    Finds a course by name or by teacher name and returns it.
* **URL**

  /courses/:by/:pattern
  
* **Method:**

  `GET`
  
*  **URL Params**

   **Required:**
 
   `by=[string]`<br/>
   `pattern=[string]`
   
   `by` can be `name` or `teacher`<br />
  `pattern` is the text that the course name or the teacher should contain to match with the search criteria 


* **Data Params**

    None

* **Success Response:**

  * **Code:** 200 OK <br />
    **Content:** [`{id:"5cda791f5aec95bb5a5abd7c",
                   name:"Advanced Calculus", department:"Science", teacher: "Doe", year: "2019-2020", semester: 2,
	               description: "Limits, Derivatives, Integrals",
	               schedule:[{day: "lun", start_time: "10:00", end_time: "11:00", room: "A4" },
				                {day: "mar", start_time: "11:30", end_time: "12:30", room: "B9"}]
                     },
                   {id:"5cda791f5aec95bb5a5abd7c",
                   name:"Machine Learning", department:"Computer Science", teacher: "Doe", year: "2019-2020", semester: 1,
	               description: "SVM, Clustering, Neural Networks",
	               schedule:[{day: "lun", start_time: "10:00", end_time: "11:00", room: "A4" },
				                {day: "mar", start_time: "11:30", end_time: "12:30", room: "B9"}]
                     }]`
 
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
    **Content:** `{ error : "Permission denied" }` This error may occur as the course creation is allowed to teachers only.
    
  OR

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{ error : "Expired token" }`