**Find Teaching Material By Course**
----
    Finds all teaching material from a specific course.
* **URL**

  /teachingMaterials/:courseId
  
* **Method:**

  `GET`
  
*  **URL Params**

   **Required:**
 
   `courseId=[string]`<br/>
   `courseId` is the id of the course which teaching material are searched <br />

* **Data Params**

    None

* **Success Response:**

  * **Code:** 200 OK <br />
      **Content:** `[
          "Temp1.xml",
          "Temp2.txt"
          ]`
 
* **Error Response:**

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
    **Content:** `{ error : "Permission denied" }` <br />
    This error may occur as the course creation is allowed to teachers only.
    
  OR

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{ error : "Expired token" }`