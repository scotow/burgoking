document.addEventListener('DOMContentLoaded', function() {
    function setLabel(text) {
        document.getElementById('label').innerText = text;
    }

    function errorOccurred() {
        const genericError = 'An error has occurred. Feel free to send me an email to the admin of the website.';

        request.onreadystatechange = function() {
            if (request.readyState === XMLHttpRequest.DONE) {
                if (request.status === 200) {
                    setLabel('An error has occurred. Feel free to send me an email at \'' + request.responseText + '\' to help me improve this project.');
                } else {
                    setLabel(genericError);
                }
            }
        };
        request.onerror = function () {
            setLabel(genericError);
        };

        request.open('GET', '/contact', true);
        request.send(null);
    }

    const request = new XMLHttpRequest();

    request.onreadystatechange = function() {
        if (request.readyState === XMLHttpRequest.DONE) {
            if (request.status === 200) {
                setLabel(request.responseText);
            } else {
                errorOccurred()
            }
        } else if (request.readyState === XMLHttpRequest.LOADING) {
            setLabel('Loading your Burger');
        }
    };
    request.onerror = errorOccurred;

    request.open('POST', '/code', true);
    request.send(null);
});