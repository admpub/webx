
// Validate Contact Form
$("#ContactForm").validate({
    rules: {
        full_name: {
            required: true,
            minlength: 3
        },
        email: {
            required: true,
            minlength: 4
        }
        ,
        mobile: {
            required: true
        },
        subject: {
            required: true
        },
        message: {
            required: true
        }
    },
    messages: {
        full_name: {
            required: "Please Enter Name",
            minlength: "Name must consist of at least 3 characters"
        },
        email: {
            required: "Please provide a Email",
            minlength: "Email must be at least 4 characters long"
        },
        mobile: {
            required: "Please provide Mobile Number"
        },
        subject: {
            required: "Please Enter Subject"
        },
        message: {
            required: "Please Enter Message"
        }
    },

    submitHandler: function(form) {
  
        var formdata = jQuery("#ContactForm").serialize();
        jQuery.ajax({
            type:"POST",
            url:"assets/contact-form/contact.php",
            data:formdata,
            dataType:'json',
            async: false,
            success:function(data){
                if(data.success){
                    alert('Thank You, Your Message Has been Sent');
                }else{
                    alert('Error on Sending Message, Please Try Again');
                }

            },
            error:function(error){
                alert('Something Went Wrong');
            }

        });
    }
});

