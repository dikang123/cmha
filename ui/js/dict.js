var IP = "192.168.200.135:8500";

var FreshenTime = 5e3;

var IP_old = 0;

var serviceName = "";

var hostName = "";

var HostType = "";

var Hostrole_id = "";

var formatter_db_status = function(cellvalue, options, rowObject) {
    if (rowObject.REPL - STATUS == "OK") {
        return '<span style="color:green;" >' + cellvalue + "</span>";
    } else if (rowObject.REPL - STATUS == "warning") {
        return '<span style="color:red;" >' + cellvalue + "</span>";
    } else {
        return '<span style="color:red;" >' + cellvalue + "</span>";
    }
};

var formatter_db_counter_status = function(cellvalue, options, rowObject) {
    if (rowObject.REPL - ERR - COUNTER == 0) {
        return '<span style="color:green;" >' + cellvalue + "</span>";
    } else {
        return '<span style="color:red;" >' + cellvalue + "</span>";
    }
};

var formatter_role = function(cellvalue, options, rowObject) {
    if (Hostrole_id == "leader") {
        return '<span  style="color:green;" >' + Hostrole_id + "</span>";
    }
    return "<span  >" + cellvalue + "</span>";
};

var formatter_repl_status = function(cellvalue, options, rowObject) {
    if (rowObject.REPL_STATUS == "OK") {
        return '<span style="color:green;" >' + cellvalue + "</span>";
    } else if (rowObject.REPL_STATUS == "warning") {
        return '<span style="color:red;" >' + cellvalue + "</span>";
    } else {
        return '<span style="color:red;" >' + cellvalue + "</span>";
    }
};

var formatter_counter_status = function(cellvalue, options, rowObject) {
    if (rowObject.REPL_ERR_COUNTER == 1) {
        return '<span style="color:red;" >' + cellvalue + "</span>";
    }
    return '<span style="color:green;" >' + cellvalue + "</span>";
};