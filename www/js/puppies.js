$(document).ready(function() {
    $('.puppy-image-container').hover(
    function () {
        $(this).find('.puppy-image-upvote, .puppy-image-downvote').fadeIn(200);
    },
    function () {
        $(this).find('.puppy-image-upvote, .puppy-image-downvote').fadeOut(500);
    });
});