var load = function() {
    $("#content").load("ajax/jqgrid.html");
};

var indexfunction = function() {
    window.location.href = "http://" + IP + "/ui";
};

var backFunction = function() {
    window.location.href = "http://" + IP + "/ui/log.html";
};

$("#Home").click(function() {
    $("#content").load("ajax/jqgrid.html");
});

$(".serviceNews").click(function() {
    serviceName = $(this).attr("id");
    document.cookie = "serviceName=" + serviceName;
    $("#content").load("ajax/jqgrid2.html");
});

$(".serviceHost").click(function() {
    hostName = $(this).attr("id");
    document.cookie = "hostName=" + hostName;
    $("#content").load("ajax/jqgrid3.html");
});