$(document).ready(function() {

    var ImageHandler = function() {
        this.init();
    };

    ImageHandler.prototype = {
        init: function () {
            this.showVotes();
            this.bindVotingEvents();
        },

        showVotes: function() {
            $('.puppy-image-container').hover(
            function () {
                $(this).find('.puppy-image-upvote, .puppy-image-downvote, .puppy-image-upvote-info, .puppy-image-downvote-info').fadeIn(200);
            },
            function () {
                $(this).find('.puppy-image-upvote, .puppy-image-downvote, .puppy-image-upvote-info, .puppy-image-downvote-info').fadeOut(500);
            });
        },

        bindVotingEvents: function() {
            var self = this;
            $('.puppy-image-upvote').unbind('click').bind('click', function() {
                var img_id = $(this).parent().attr("img_id");
                self.votePuppy(img_id, true);
            });
            $('.puppy-image-downvote').unbind('click').bind('click', function() {
                var img_id = $(this).parent().attr("img_id");
                self.votePuppy(img_id, false);
            });
        },

        votePuppy: function (img_id, is_upvote) {
            var self = this;
            $.ajax({
                url: 'vote',
                data: {
                    img_id: img_id,
                    is_upvote: is_upvote
                },
                type: 'POST',
                success: function(new_vote) {
                    if(is_upvote) {
                        $('.puppy-image-container[img_id="' + img_id + '"]').find('.puppy-image-upvote-info').html(new_vote);
                    } else {
                        $('.puppy-image-container[img_id="' + img_id + '"]').find('.puppy-image-downvote-info').html(new_vote);
                    }
                }
            });
        }
    };

    new ImageHandler();
});