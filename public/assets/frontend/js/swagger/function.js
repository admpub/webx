function requestInterceptor() {
    // 生成签名
    // @author Hank Shen <swh@admpub.com>
    var docURL = $('input.download-url-input').val();
    if (!docURL||this.url == 'doc.json'||this.url == docURL) return this; // 打开页面时加载doc.json
    this.headers['X-App-ID'] = $('#app-id-input').val();
    this.headers['X-App-Sign'] = '';
    this.headers['X-App-Timestamp'] = new Date().getTime();
    var pos = this.url.indexOf('?'), params = 'appID=' + this.headers['X-App-ID'];
    var body = '';
    if (typeof (this.body) != 'undefined') body = this.body;
    if (pos >= 0) {
        params = this.url.substring(pos) + '&' + params;
    } else {
        params = '?' + params;
    }
    var curURL = window.location.href;
    var url = curURL.substring(0, curURL.lastIndexOf('/')) + '/makesign' + params;
    var that = this;
    $.ajax({
        url: url,
        type: 'POST',
        data: body,
        dataType: 'json',
        async: false,
        headers: this.headers,
        success: function (r) {
            if (r.Code != 1) return alert(r.Info);
            that.headers['X-App-Sign'] = r.Data.sign;
        }, error: function (r) {
            if (r.responseJSON) {
                if (r.responseJSON.Info) alert(r.responseJSON.Info);
                else alert(r.responseText);
            } else alert(r.responseText);
        }
    });
    //console.dir(this); //DEBUG
    return this;
}

function onSwaggerReady() {
}