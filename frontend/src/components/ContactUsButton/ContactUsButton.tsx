// node modules
import React from 'react';

type ContactUsButtonProps = {
  children: any;
  buttonElem: Function;
};

const ContactUsButton: React.FC<ContactUsButtonProps> = ({
  children,
  buttonElem
}) => {
  const onConnectButtonClick = () => {
    const anchor = document.querySelector('#contact-form-anchor');

    if (anchor) {
      anchor.scrollIntoView({
        behavior: 'smooth',
        block: 'center',
      });
    }
  }

  const Button = buttonElem();
  console.log(Button);
  

  return (
    <Button onClick={() => onConnectButtonClick()}>
      { children }
    </Button>
  );
}

export default ContactUsButton;
