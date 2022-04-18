if(typeof String.prototype.trim !== 'function') {
    String.prototype.trim = function() {
        return this.replace(/^\s+|\s+$/g, ''); 
    }
}

var PVCC = PVCC || {
    setCookie: function(cname,cvalue,exsecs)
    {   
        var d = new Date();
        d.setTime(d.getTime()+(exsecs*1000));
        var expires = "expires="+d.toGMTString();
        document.cookie = cname + "=" + cvalue + ";path=/; " + expires;
    },  
    getCookie: function(cname)
    {   
        var name = cname + "=";
        var ca = document.cookie.split(';');
        for(var i=0; i<ca.length; i++) 
        {   
            var c = ca[i].trim();
            if (c.indexOf(name)==0) return c.substring(name.length,c.length);
        }   
        return ""; 
    },
    getVar: function(name) {
        name = name.replace(/[\[]/, "\\[").replace(/[\]]/, "\\]");
        var regex = new RegExp("[\\?&]" + name + "=([^&#]*)"), results = regex.exec(location.search);
        return results === null ? "" : decodeURIComponent(results[1].replace(/\+/g, " "));
    }
};

function selectLang(lang){
    PVCC.setCookie('sodht_language', lang, 86400);
    location.reload();
}

function goLink(url){
    var get_category = getQueryString('category');
    if (get_category){
        location.href = url+'?category='+get_category;
    } else {
        location.href = url;
    }
}

function dataSortBy(obj, stype, num) {
    var tpl = $(obj).attr('tpl');

    var type_list = [0,1,2];
    if (type_list.indexOf(num) != -1 ){
        if (num + 1 <=2){
            num++;
        } else {
            num = 0;
        }
    } else {
        num = 0;
    }


    tpl = tpl.replace('{num}', num);
    location.href = tpl;
}