$(document).ready(function() {

    var ImageHandler = function() {
        this.init();
    };

    ImageHandler.prototype = {
        init: function () {
            this.showVotes();
            this.bindVotingEvents();
            this.bindSort();
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
        },

        bindSort: function() {
            var self = this;
            $('.sort-btn').unbind('click').bind('click', function() {
                var sort_type = $('.js-sort-type');
                if(sort_type.attr('type') == "asc") {
                    var data_sort_type = "desc";
                    sort_type.html("(DESC)");
                    sort_type.attr("type", "desc");
                } else {
                    data_sort_type = "asc";
                    sort_type.html("(ASC)");
                    sort_type.attr("type", "asc");
                }
                var image_collection = [];
                var puppy_images = $('.puppy-image-container');
                $.each(puppy_images, function(index, item) {
                    image_collection.push({
                        img_id: $(item).attr("img_id"),
                        up_vote: $(item).find('.puppy-image-upvote-info').html()
                    });
                });
                var sorted_images = image_collection.sort(function(obj1, obj2) { return obj1.up_vote - obj2.up_vote });
                if(data_sort_type == "asc") {
                    sorted_images = sorted_images.reverse();
                }

                var sorted_image_collection = [];
                $.each(sorted_images, function(index, item) {
                    var img = $('.puppy-image-container[img_id="' + item.img_id + '"]');
                    sorted_image_collection.push(img);
                });


                $('.body-container').empty();
                $.each(sorted_image_collection, function(i, item) {
                    $('.body-container').append(item);
                });

                self.showVotes();
                self.bindVotingEvents();
            });
        }
    };

    new ImageHandler();
});