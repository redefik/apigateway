**Get Student's Courses**
----
    Returns the courses which a student is subscribed to.
* **URL**

  /courses/students/:student_username
  
* **Method:**

  `GET`
  
*  **URL Params**

   **Required:**
 
   `student_username=[string]`

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
    **Content:** `{ error : "Courses Not Found"}`
    This is returned where no courses are associated to the student
  
  OR

  * **Code:** 404 NOT FOUND <br />
    **Content:** `{ error : "Student Not Found" }`
    This is returned where no student with the given username exist
    
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