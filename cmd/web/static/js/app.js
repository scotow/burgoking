document.addEventListener('DOMContentLoaded', function() {
    function errorOccurred() {
        document.getElementById('label').innerText = 'An error has occurred. Feel free to send me an email at \'contact@scotow.com\' to help me improve this project.';
    }

    const request = new XMLHttpRequest();

    request.onreadystatechange = function() {
        if (request.readyState === XMLHttpRequest.DONE) {
            if (request.status === 200) {
                document.getElementById('label').innerText = request.responseText;
            } else if (request.status === 400) {
                errorOccurred()
            }
        }
    };
    request.onerror = errorOccurred;

    request.open('GET', '/code', true);
    request.send();

    document.getElementById('label').innerText = 'Loading your Burger';
});