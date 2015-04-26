package main

const Template = `<!doctype html>
<html>
<head>
    <title>Wiki</title>
    <meta charset="utf-8" />

    <style type="text/css">
        body {
            margin: 1rem;
            font-family: georgia;
            font-size: 0.9rem;
        }

        pre {
            font-size: 1rem;
            margin: 0;
            padding 0;
        }

        ul {
            list-style-type: square;
            list-style-position: inside;
            padding-left: 1rem;
        }

        hr {
            height: 1px;
            border: none;
            border-bottom: 1px #eee solid;
        }

        h1,
        h2,
        h3,
        h4,
        h5,
        h6 {
            font-family: 'Helvetica Neue', georgia;
            font-weight: 200;
        }

        .commit {
            cursor: pointer;
            background-color: #eee;
            margin: 0.2rem 0;
            padding: 0.2rem;
        }

        .highlight .gi > .x {
            background-color: #eaffea;
            color: #55a532;
        }

        .highlight .gd > .x {
            background-color: #ffecec;
            color: #bd2c00;
        }

        .highlight {
            margin: 0.5rem 0;
            border-left: 3px solid #eee;
            padding-left: 0.5rem;
        }

        .date,
        .hash,
        .author,
        .subject {
            display: inline-block;
            text-overflow: ellipsis;
            margin: 0.1rem 0;
        }
    </style>

    {{ if .CustomCSS }}
    <link rel="stylesheet" href="/css/custom.css" type="text/css">
    {{ end }}

</head>
<body>
    <header>{{ .Title }}</header>

    <section>{{ .Body }}</section>

    <hr>

    <footer class="commits">
        {{ range .Commits }}
        <div class="commit" data-hash="{{ .Hash }}" data-file="{{ .FileNoExt }}">
            <span class="date">{{ .HumanDate }}</span> &middot;
            <span class="author">{{ .Author }}</span> &middot;
            <span class="subject">{{ .Subject }}</span>
            <div class="diff"></div>
        </div>
        {{ end}}
    </footer>

    <script>
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
    </script>
</body>
</html>`
