
$(function(){
	$("#slidey").slidey({
		interval: 8000,
		listCount: 5,
		autoplay: false,
		showList: true
	});
	$(".slidey-list-description").dotdotdot();

	$("#owl-demo").owlCarousel({
	  autoPlay: 3000, //Set AutoPlay to 3 seconds
	  items : 5,
	  itemsDesktop : [640,4],
	  itemsDesktopSmall : [414,3]
	});

	$('.sign-in-modal .toggle').click(function(){
		// Switches the Icon
		$(this).children('i').toggleClass('fa-pencil');
		// Switches the forms  
		$('.sign-in-modal .form').animate({
			height: "toggle",
			'padding-top': 'toggle',
			'padding-bottom': 'toggle',
			opacity: "toggle"
		}, "slow");
	});

	$('.flexslider').flexslider({
		animation: "slide",
		start: function(slider){
			$('body').removeClass('loading');
		}
	});
});