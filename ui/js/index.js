var changeMemu = function() {
    var dataServices = {};
    var after_dataMemu = [];
    var dataMenuItem = {};
    var dataAllService = {};
    var getService_db = function() {
        $.ajax({
            url:"http://" + IP + "/v1/catalog/services",
            method:"get",
            async:false,
            dataType:"json",
            success:function(result, status, xhr) {
                dataServices = result;
            },
            error:function(XMLHttpRequest, status, jqXHR, textStatus, e) {
                console.error("失败状态文本 " + status);
            }
        });
    };
    var getAllDataService = function() {
        for (var i = 0; i < after_dataMemu.length; i++) {
            $.ajax({
                method:"get",
                url:"http://" + IP + "/v1/health/service/" + after_dataMemu[i],
                async:false,
                dataType:"json",
                success:function(result, status, xhr) {
                    var cmha_data_service_old = [];
                    cmha_data_service_old = result;
                    dataAllService[after_dataMemu[i]] = cmha_data_service_old;
                },
                error:function(XMLHttpRequest, status, jqXHR, textStatus, e) {
                    console.error("失败状态文本 " + status);
                }
            });
        }
    };
    var changeDataAllSrvice = function() {
        try {
            for (var ky in dataAllService) {
                var serviceName = [];
                for (var a = 0; a < dataAllService[ky].length; a++) {
                    serviceName.push(dataAllService[ky][a].Node.Node);
                }
                dataMenuItem[ky] = serviceName;
            }
        } catch (erro) {
            console.error("changeDataAllSrvice  获得主机服务名出错！！！");
        }
    };
    Array.prototype.indexOf = function(val) {
        for (var i = 0; i < this.length; i++) {
            if (this[i] == val) return i;
        }
        return -1;
    };
    Array.prototype.remove = function(val) {
        var index = this.indexOf(val);
        if (index > -1) {
            this.splice(index, 1);
        }
    };
    var changeServiceName = function() {
        try {
            for (var key in dataServices) {
                after_dataMemu.push(key);
            }
            after_dataMemu.remove("Statistics");
            after_dataMemu.remove("consul");
        } catch (erro) {
            console.error("changeServiceName 筛选服务名出错 ！！！");
        }
    };
    function addElementLi(obj) {
        var ul = document.getElementById(obj);
        for (var k in dataMenuItem) {
            var li = document.createElement("li");
            var a = document.createElement("a");
            a.setAttribute("href", "#");
            var i = document.createElement("i");
            i.setAttribute("class", "fa fa-skype fa-lg");
            var ul_1 = document.createElement("ul");
            var li_news = document.createElement("li");
            var a_news = document.createElement("a");
            a_news.setAttribute("href", "#");
            var i_news = document.createElement("i");
            i_news.setAttribute("class", "fa fa-stack-exchange fa-lg ");
            var span_news = document.createElement("span");
            span_news.setAttribute("class", "serviceNews nav-text");
            span_news.setAttribute("id", k);
            span_news.innerHTML = " base info";
            a_news.appendChild(i_news);
            a_news.appendChild(span_news);
            li_news.appendChild(a_news);
            ul_1.appendChild(li_news);
            for (var ky = 0; ky < dataMenuItem[k].length; ky++) {
                var li_1 = document.createElement("li");
                var a_1 = document.createElement("a");
                a_1.setAttribute("href", "#");
                var i_1 = document.createElement("i");
                i_1.setAttribute("class", "fa fa-dot-circle-o fa-lg  ");
                var span_1 = document.createElement("span");
                span_1.setAttribute("class", "serviceHost nav-text");
                span_1.setAttribute("id", dataMenuItem[k][ky]);
                span_1.innerHTML = dataMenuItem[k][ky];
                a_1.appendChild(i_1);
                a_1.appendChild(span_1);
                li_1.appendChild(a_1);
                ul_1.appendChild(li_1);
            }
            var span = document.createElement("span");
            span.setAttribute("class", "nav-text");
            span.innerHTML = k;
            a.setAttribute("href", "#");
            a.appendChild(i);
            a.appendChild(span);
            li.appendChild(a);
            li.appendChild(ul_1);
            ul.appendChild(li);
        }
    }
    var startFunction = function() {
        getService_db();
        changeServiceName();
        getAllDataService();
        changeDataAllSrvice();
        addElementLi("parentUl");
    };
    startFunction();
};

changeMemu();