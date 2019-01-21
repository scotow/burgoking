document.addEventListener('DOMContentLoaded', function() {
    function setLabel(text) {
        document.getElementById('label').innerText = text;
    }

    function sendRequest(method, url, onLoad, onSuccess, onError) {
        const request = new XMLHttpRequest();

        request.onreadystatechange = function() {
            if (request.readyState === XMLHttpRequest.DONE) {
                if (request.status === 200) {
                    if (onSuccess) onSuccess(request);
                } else {
                    if (onError) onError(request);
                }
            } else if (request.readyState === XMLHttpRequest.OPENED) {
                if (onLoad) onLoad();
            }
        };
        request.onerror = onError;

        request.open(method, url, true);
        request.send(null);
    }

    function errorOccurred() {
        sendRequest('GET', '/contact', null, function (request) {
            setLabel('An error has occurred. Feel free to send me an email at \'' + request.responseText + '\' to help me improve this project.');
        }, function () {
            setLabel('An error has occurred. Feel free to send an email to the admin of the website.');
        });
    }

    sendRequest('POST', '/code', function() {
        setLabel('Loading your Burger');
    }, function (request) {
        setLabel(request.responseText);
        document.title = request.responseText;
        new Audio('/sounds/burger.m4a').play();
    }, errorOccurred);
});