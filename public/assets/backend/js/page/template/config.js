(function(){
function initSelectpage(){
    $('input.form-selectpage').each(function(){
        var multiple=$(this).data('multiple')?true:false;
        var maxSelectLimit=$(this).data('limit');
        if(maxSelectLimit) {
            maxSelectLimit = Number(maxSelectLimit);
            if(!multiple && maxSelectLimit>1) multiple = true;
        }
        var url=$(this).data('url');
        if(!url && (url=$(this).data('json'))){
            if(typeof url == 'string') url = JSON.parse(url);
        }
        if(!url) url=BACKEND_URL + '/official/article/index';
        var showField=$(this).data('showfield')||'title';
        var keyField=$(this).data('keyfield')||'id';
        var thumbField=$(this).data('thumbfield')||'image';
        $(this).selectPage({
            showField: showField,
            keyField: keyField,
            data: url,
            multiple : multiple,
            maxSelectLimit : maxSelectLimit,
            params: function () { return {}; },
            formatItem: function (row) {
                if(row[thumbField]) return '<img src="'+row[thumbField]+'" />'+row[showField];
                return row[showField];
            },
            eAjaxSuccess: function (d) {
                var result;
                if (d && d.Data && d.Data.listData && d.Data.pagination) result = {
                    "list": d.Data.listData,
                    "pageSize": d.Data.pagination.limit,
                    "pageNumber": d.Data.pagination.page,
                    "totalRow": d.Data.pagination.rows,
                    "totalPage": d.Data.pagination.pages
                };
                else result = undefined;
                return result;
            }
        });
    })
}
$(function () {
    initSelectpage();
});
})();