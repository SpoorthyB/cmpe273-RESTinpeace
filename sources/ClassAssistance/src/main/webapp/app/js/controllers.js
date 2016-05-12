
var classAssistance = angular.module('classAssistance', []);

classAssistance.controller('ClassAssistanceController', ['$scope', '$filter',
    function ($scope, $filter) {

        $scope.studentInformation = [];
        $scope.oldTable = [];
        //$scope.public_ip = 'localhost';
        $scope.public_ip = '52.37.72.212:3005';
        $scope.showAddMessage = "";
        $scope.showPresentMessage = "";
        $scope.showAbsentMessage = "";
        
        $scope.init = function () {
            //initCharts();
            $scope.getAllStudents();
            //$scope.getNumberSensorsByTime($scope.date);
        };

        $scope.getAllStudents = function () {
            $scope.showAddMessage = "";
            $scope.showPresentMessage = "";
            $scope.showAbsentMessage = "";
            var candidate = 'trump';
            var xhr = new XMLHttpRequest();
            var url = "http://" + $scope.public_ip;
            //var url = "http://54.201.81.197:8080/SentimentAnalysis/api/admin/getNumberTweetsByHashtagAndSentiment/";
            var xhr = new XMLHttpRequest();
            xhr.open("GET", url, false);
            //xhr.setRequestHeader("Content-type", "application/json");            
            xhr.send();

            if (xhr.status === 200) {
                var allUsersJSON = JSON.parse(xhr.responseText);

                $scope.studentInformation.splice(0, $scope.studentInformation.length);                
                for (i = 0; i < allUsersJSON.students.length; i++) {
                    $scope.studentInformation[i] = [];
                    $scope.studentInformation[i]["id"] = allUsersJSON.students[i].id;
                    $scope.studentInformation[i]["name"] = allUsersJSON.students[i].name;
                    $scope.studentInformation[i]["status"] = allUsersJSON.students[i].attended;
                    $scope.studentInformation[i]["date"] = allUsersJSON.students[i].time_stamp;

                }
                for (i = 0; i < $scope.studentInformation.length; i++) {
                    var bol = false
                    
                    if (i >= $scope.oldTable.length){
                        bol = true
                    }else if ($scope.studentInformation[i]["status"] != $scope.oldTable[i]["status"]){
                        bol = true
                    }
                                  
                    if (bol == true) {
                        if ($scope.studentInformation[i]["status"] == "yes") {
                            $scope.showPresentMessage = "The user " + $scope.studentInformation[i]["name"] + " is present";
                        } else if ($scope.studentInformation[i]["status"] == "no") {
                            $scope.showAbsentMessage = "The user " + $scope.studentInformation[i]["name"] + " is absent";
                        } else {
                            $scope.showAddMessage = "The user " + $scope.studentInformation[i]["name"] + " has been registered";
                        }
                    }
                }
                $scope.oldTable = [];
                for (i = 0; i < $scope.studentInformation.length; i++) {
                    $scope.oldTable[i] = [];                    
                    $scope.oldTable[i]["status"] = $scope.studentInformation[i]["status"];                  
                }
                

                //$scope.$apply();
            }


        };


        setInterval(function () {
            $scope.getAllStudents();
            $scope.$apply()
        }, 3000);

    }]);


