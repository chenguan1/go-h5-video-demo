var messageTxt = document.getElementById("messageTxt");
var messages = document.getElementById("messages");
var sendBtn = document.getElementById("sendBtn")

w = new WebSocket("ws://" + HOST + "/my_endpoint");
w.onopen = function () {
    console.log("Websocket connection enstablished");
};

w.onclose = function () {
    appendMessage("<div><center><h3>Disconnected</h3></center></div>");
};
w.onmessage = function (message) {
    appendMessage("<div>" + message.data + "</div>");

    sg = message.data.split(';')
    url = sg[sg.length-1]
    urlToImage(url)


};

sendBtn.onclick = function () {
    myText = messageTxt.value;
    messageTxt.value = "";

    appendMessage("<div style='color: red'> me: " + myText + "</div>");
    w.send(myText);
};

messageTxt.addEventListener("keyup", function (e) {
    if (e.keyCode === 13) {
        e.preventDefault();

        sendBtn.click();
    }
});

function appendMessage(messageDivHTML) {
    messages.insertAdjacentHTML('afterbegin', messageDivHTML);
}


function urlToImage(str12) {
    url = 'data:image/png;base64,'+str12;

    var outputImg = document.getElementsByTagName("img")[0];
    outputImg.src = url;

    /*var drawing = document.getElementById('drawing');
    if(drawing.getContext) {
        var context = drawing.getContext('2d');
        var slider = document.getElementById('scale-range');
        var W = 400;
        var H = 300;
        drawing.width = W;
        drawing.height = H;
        var image = new Image();
        image.src = url;
        image.onload = function(){
            context.clearRect(0,0,W,H);
            context.drawImage(image,0,0,W.H);
        }
    }*/
}



