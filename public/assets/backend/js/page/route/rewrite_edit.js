function searchRoute(){
    var size = 10;
    $('#input-route').typeahead({hint: true, highlight: true, minLength: 1}, {
        source: function (query, sync, async) {
            var data = { prefix: query, size: size };
            $.ajax({
                url: '?op=routeList',
                type: 'get',
                data: data,
                dataType: 'json',
                async: false,
                success: function (data) {
                    var arr = [];
                    if (!data.Data) return;
                    $.each(data.Data, function (index, val) {
                        arr.push(val);
                    });
                    sync(arr);
                }
            });
        }, limit: size
    });
}
$(function(){
    searchRoute()
})