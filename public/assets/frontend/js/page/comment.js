$(function() {
    initPageAllCommentForm();
    var $container = $('#comment-list-container');
    if($container.length > 0 && $container.has('[scroll-load]')) {
        $container.scrollLoading({callback:function(){commentList(1);},attr:''});//attr:设置为空避免自动加载，默认为data-url
    }
});