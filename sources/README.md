#CMPE 273 - Class project

##The system implements an automated attendence supervising system. It uses the BLE components of smartphones and RaspberryPi to check for the attendance.

##The system architecture is as follows:
###Android App - It is used to register the student details for the class. It sends frquent Bluetooth advertisements that are detected by the raspberryPi
###RaspberryPI - It receives advertisements from the smartphone and transmits the data to the backend server.
###Backend server - It handles the REST calls made by the smartphone and the raspberryPi and updates the student information.
###Frontend server - It displays class attendance by showing student details in real-time.

 
 
