<<<<<<< HEAD
; (function($, window, document, undefined) {
    var pluginName = "jqueryAccordionMenu";
    var defaults = {
        speed: 300,
        showDelay: 0,
        hideDelay: 0,
        singleOpen: true,
        clickEffect: true
    };
    function Plugin(element, options) {
        this.element = element;
        this.settings = $.extend({},
        defaults, options);
        this._defaults = defaults;
        this._name = pluginName;
        this.init()
    };
    $.extend(Plugin.prototype, {
        init: function() {
            this.openSubmenu();
            this.submenuIndicators();
            if (defaults.clickEffect) {
                this.addClickEffect()
            }
        },
        openSubmenu: function() {
            $(this.element).children("ul").find("li").bind("click touchstart",
            function(e) {
                e.stopPropagation();
                e.preventDefault();
                if ($(this).children(".submenu").length > 0) {
                    if ($(this).children(".submenu").css("display") == "none") {
                        $(this).children(".submenu").delay(defaults.showDelay).slideDown(defaults.speed);
                        $(this).children(".submenu").siblings("a").addClass("submenu-indicator-minus");
                        if (defaults.singleOpen) {
                            $(this).siblings().children(".submenu").slideUp(defaults.speed);
                            $(this).siblings().children(".submenu").siblings("a").removeClass("submenu-indicator-minus")
                        }
                        return false
                    } else {
                        $(this).children(".submenu").delay(defaults.hideDelay).slideUp(defaults.speed)
                    }
                    if ($(this).children(".submenu").siblings("a").hasClass("submenu-indicator-minus")) {
                        $(this).children(".submenu").siblings("a").removeClass("submenu-indicator-minus")
                    }
                }
                window.location.href = $(this).children("a").attr("href")
            })
        },
        submenuIndicators: function() {
            if ($(this.element).find(".submenu").length > 0) {
                $(this.element).find(".submenu").siblings("a").append("<span class='submenu-indicator'>+</span>")
            }
        },
        addClickEffect: function() {
            var ink, d, x, y;
            $(this.element).find("a").bind("click touchstart",
            function(e) {
                $(".ink").remove();
                if ($(this).children(".ink").length === 0) {
                    $(this).prepend("<span class='ink'></span>")
                }
                ink = $(this).find(".ink");
                ink.removeClass("animate-ink");
                if (!ink.height() && !ink.width()) {
                    d = Math.max($(this).outerWidth(), $(this).outerHeight());
                    ink.css({
                        height: d,
                        width: d
                    })
                }
                x = e.pageX - $(this).offset().left - ink.width() / 2;
                y = e.pageY - $(this).offset().top - ink.height() / 2;
                ink.css({
                    top: y + 'px',
                    left: x + 'px'
                }).addClass("animate-ink")
            })
        }
    });
    $.fn[pluginName] = function(options) {
        this.each(function() {
            if (!$.data(this, "plugin_" + pluginName)) {
                $.data(this, "plugin_" + pluginName, new Plugin(this, options))
            }
        });
        return this
    }
})(jQuery, window, document);
=======
!function(i,n,e,s){function t(n,e){this.element=n,this.settings=i.extend({},u,e),this._defaults=u,this._name=h,this.init()}var h="jqueryAccordionMenu",u={speed:300,showDelay:0,hideDelay:0,singleOpen:!0,clickEffect:!0};i.extend(t.prototype,{init:function(){this.openSubmenu(),this.submenuIndicators(),u.clickEffect&&this.addClickEffect()},openSubmenu:function(){i(this.element).children("ul").find("li").bind("click touchstart",function(e){if(e.stopPropagation(),e.preventDefault(),i(this).children(".submenu").length>0){if("none"==i(this).children(".submenu").css("display"))return i(this).children(".submenu").delay(u.showDelay).slideDown(u.speed),i(this).children(".submenu").siblings("a").addClass("submenu-indicator-minus"),u.singleOpen&&(i(this).siblings().children(".submenu").slideUp(u.speed),i(this).siblings().children(".submenu").siblings("a").removeClass("submenu-indicator-minus")),!1;i(this).children(".submenu").delay(u.hideDelay).slideUp(u.speed),i(this).children(".submenu").siblings("a").hasClass("submenu-indicator-minus")&&i(this).children(".submenu").siblings("a").removeClass("submenu-indicator-minus")}n.location.href=i(this).children("a").attr("href")})},submenuIndicators:function(){i(this.element).find(".submenu").length>0&&i(this.element).find(".submenu").siblings("a").append("<span class='submenu-indicator'>+</span>")},addClickEffect:function(){var n,e,s,t;i(this.element).find("a").bind("click touchstart",function(h){i(".ink").remove(),0===i(this).children(".ink").length&&i(this).prepend("<span class='ink'></span>"),n=i(this).find(".ink"),n.removeClass("animate-ink"),n.height()||n.width()||(e=Math.max(i(this).outerWidth(),i(this).outerHeight()),n.css({height:e,width:e})),s=h.pageX-i(this).offset().left-n.width()/2,t=h.pageY-i(this).offset().top-n.height()/2,n.css({top:t+"px",left:s+"px"}).addClass("animate-ink")})}}),i.fn[h]=function(n){return this.each(function(){i.data(this,"plugin_"+h)||i.data(this,"plugin_"+h,new t(this,n))}),this}}(jQuery,window,document);
>>>>>>> 228e7f9e1916a71f7393da1e52adce215fb51100
