$(function() {
    initPageAllCommentForm();
    var $container = $('#comment-list-container');
    if($container.length > 0 && $container.hasAttr('scroll-load')) {
        $container.scrollLoading({callback:function() {commentList(1);},attr:'data-url'});
    }
});