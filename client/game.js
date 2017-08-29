var w = window;


function updateTime() {
    if (prop.state.timeLeft <= 0) {
        return;
    }

    prop.state.timeLeft--;
    var minutes = Math.floor(prop.state.timeLeft / 60);
    var seconds = prop.state.timeLeft % 60;
    var min = "" + minutes;
    var ppp = "00";
    var mmm = ppp.substring(0, ppp.length - min.length) + min;
    var sec = "" + seconds;
    var sss = ppp.substring(0, ppp.length - sec.length) + sec;
    document.getElementById('timer').innerHTML = '' + mmm + ':' + sss;
}

function updateElixir() {
    if (prop.state.elixir.quantity >= 30) {
        return
    }

    prop.state.elixir.quantity++;
    document.getElementById('elixir').innerHTML = '' + Math.floor(prop.state.elixir.quantity/3)
}

function updateButtons() {
    document.getElementById('carta_formica').innerHTML = prop.conf.carte.formica.name + ' ' + prop.conf.carte.formica.elixir
    document.getElementById('carta_calabrone').innerHTML = prop.conf.carte.calabrone.name + ' ' + prop.conf.carte.calabrone.elixir
}

var tick = function(){
    setTimeout(tick, 1000);
    updateTime();
    updateElixir();
    updateButtons();
};

w.onload = function () {
    console.log('everything starts');
    tick();

    var conn;
    var d = document;
    var msg = d.getElementById("msg");
    var log = d.getElementById("log");

    d.getElementById("btn_good_luck").onclick = function() {
        conn.send(JSON.stringify({ type: "message", message: "Buona fortuna!" }));
    }

    d.getElementById("btn_good_game").onclick = function() {
        conn.send(JSON.stringify({ type: "message", message: "Bella giocata!!" }));
    }

    d.getElementById("btn_wow").onclick       = function() {
        conn.send(JSON.stringify({ type: "message", message: "Wow!" }));
    }

    d.getElementById("carta_formica").onclick = function() {
        if (prop.state.elixir.quantity >= prop.conf.carte.formica.elixir * 3) {
            conn.send(JSON.stringify({
                type: "card",
                card: prop.conf.carte.formica
            }));
            prop.state.elixir.quantity -= prop.conf.carte.formica.elixir * 3;
            updateElixir();
        } else {
            conn.send(JSON.stringify({
                type: "message",
                message: "elixir insufficiente"
            }));
        }
    }

    d.getElementById("carta_calabrone").onclick = function() {
        if (prop.state.elixir.quantity >= prop.conf.carte.calabrone.elixir * 3) {
            conn.send(JSON.stringify({
                type: "card",
                card: prop.conf.carte.calabrone
            }));
            prop.state.elixir.quantity -= prop.conf.carte.calabrone.elixir * 3;
            updateElixir();
        } else {
            conn.send(JSON.stringify({
                type: "message",
                message: "elixir insufficiente"
            }));
        }
    }


    function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    if (w["WebSocket"]) {
        conn = new WebSocket("ws://" + d.location.host + "/ws");
        //conn = new WebSocket("ws://localhost:8080/ws");

        conn.onclose = function (evt) {
            var item = d.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        };

        conn.onmessage = function (evt) {
            var messages = evt.data.split('\n');
            var communications = d.getElementById('communications');

            json = eval(messages);
            request = JSON.parse(json);

            if (request.type == 'message') {
                communications.innerHTML = request.message;
            }

            if (request.type == 'card') {
                var item = d.createElement("div");
                item.innerText = request.card.name
                appendLog(item);
                communications.innerHTML = 'nuova carta in gioco';
            }
        };
    } else {
        var item = d.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }
};
