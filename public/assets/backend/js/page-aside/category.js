(function ($) {
    function lazyloadCategory(that, top, jumpTo, result) {
        jumpTo = String(jumpTo);
        var q = /"/g;
        return function (r) {
            if (r.Code != 1) return App.message({ text: r.Info, type: 'error' });
            var lis = '';
            for (var i = 0; i < r.Data.listData.length; i++) {
                // id,name,has_child
                var v = result?result(r.Data.listData[i]):r.Data.listData[i];
                var a = '<a href="' + jumpTo.replace('{id}',v.id) + '" data-id="' + v.id + '" data-pjax="#pcont" data-keepjs="true" title="'+String(v.name).replace(q,'&quot;')+'">' + v.name + '</a>';
                var s = '<li>';
                if (v.has_child == 'Y') {
                    s += '<label class="nav-header">';
                    s += '<i class="fa fa-folder-o tree-toggler"></i>';
                    s += a;
                    s += '</label>'
                    s += '<ul class="nav nav-list tree"></ul>'
                } else {
                    s += '<label class="nav-header">';
                    s += '<i class="fa fa-folder-open-o"></i>';
                    s += a;
                    s += '</label>'
                }
                s += '</li>';
                lis += s;
            }
            if (top) {
                $(that).append(lis);
            } else {
                $(that).parent().next('ul.tree').append(lis);
            }
        }
    }
    function cateClick(a) {
        $('.treeview').find('li.current').removeClass('current');
        $(a).closest('li').addClass('current');
    }
    function init(container, options){
        var href = $(this).data('href')||options.href, url = $(this).data('url')||options.url, result = options.result||null;
        $.get(url, {}, lazyloadCategory(container, true, href, result), 'json');
        App.treeToggle(null, {
            ajax: function () {
                var that = this;
                return {
                    url: url,
                    method: 'GET',
                    dataType: 'json',
                    data: function () {
                        return {
                            parentId: $(that).next('a').data('id') || '',
                        };
                    },
                    success: lazyloadCategory(that, false, href, result)
                }
            },
        });
    }
    /*
    options {href:"",url:"",result:function(row){ return {id:row,name:row.name,has_child:row.has_child}; }}
    */
    $.fn.pageAsideCategory = function (options) {
        $(this).each(function(){
            init(this,options);
        });
        $(this).on('click','a', function(){
            cateClick(this);
        });
    };
})(jQuery);