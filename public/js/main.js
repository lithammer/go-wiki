(function() {
    var commits = document.querySelectorAll('.commit');

    function closeDiffs(skipElement) {
        var diffs = document.querySelectorAll('.diff');
        for (var i = 0; i < diffs.length; i++) {
            if (skipElement && skipElement === diffs[i]) {
                continue;
            }

            diffs[i].innerHTML = '';
        }
    }

    function request(dataset, callback) {
        var url = '/api/diff/' + dataset.hash + '/' + dataset.file;
        var r = new XMLHttpRequest();

        r.onreadystatechange = function() {
            if (r.readyState === 4 && r.status === 200) {
                callback(r.responseText);
            }
        };

        r.open('GET', url, true);
        r.send(null);
    }

    function onClick() {
        var diff = this.querySelector('.diff');

        if (diff.innerHTML === "") {
            request(this.dataset, function(response) {
                diff.innerHTML = response;
            });
        } else {
            closeDiffs();
        }

        closeDiffs(diff);
    }

    for (var i = 0; i < commits.length; i++) {
        commits[i].addEventListener('click', onClick);
    }
})();
