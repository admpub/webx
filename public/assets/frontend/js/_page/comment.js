$(function() {
    initPageAllCommentForm();
    $('#comment-list-container').scrollLoading({
        callback: function() {
            commentList(1);
        },
        attr: '',
    });
});