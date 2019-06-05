**Reserve Exam**
----
 Registers a student to an exam
* **URL**

  /exams/:exam_id/students/:student_username

* **Method:**

  `PUT`
  
*  **URL Params**

   **Required:**
 
   `exam_id=[string]`<br/>
   `student_username=[string]`
   
   
* **Data Params**

    None

* **Success Response:**

  * **Code:** 200 OK <br />
    **Content:** `{
    "call": 2,
    "course": "IdCourse",
    "date": "21-03-2019",
    "expirationDate": "18-03-2019",
    "id": "5ce28417b8a5677e75af4288",
    "room": "A1",
    "startTime": "10:30",
    "students": [
        "5ce2a7945755bfabdf1a2aca"
    ]
}` (This is the updated exam with the identifiers of the registered students)
 
* **Error Response:**

  * **Code:** 409 CONFLICT <br />
    **Content:** `{ error : "Conflict - The resource already exists"}`
    This is returned when the student is already registered
    
  OR

  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Internal Server Error" }`
    
  OR

  * **Code:** 404 NOT FOUND <br />
    **Content:** `{ error : "Student Not Found" }`
    This is returned where a student with the given username does not exist
    
  OR

  * **Code:** 404 NOT FOUND <br />
    **Content:** `{ error : "Exam Not Found" }`
    This is returned where an exam with the given id does not exist
    
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
    
  OR

  * **Code:** 403 FORBIDDEN <br />
    **Content:** `{ error : "Not Registered To Course" }`
    This is returned when the student is not registered to the course of the exam
     