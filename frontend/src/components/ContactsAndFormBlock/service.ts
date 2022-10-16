// node modules
import { send } from 'emailjs-com';
// store
import ContactsBlockStore from './store';

class ContactsBlockService extends ContactsBlockStore {
  sendEmail = () => {
    const form = {
      name: this.name,
      description: this.description,
      phone: this.phone,
      email: this.email,
    };

    send('service_t92asb2', 'template_rfkanrl', form, 'ZS1c3SKSTN5IRcFmW')
      .then((result) => {
          console.log(result.text);
      }, (error) => {
          console.log(error.text);
      });
  }

  handleValidateForm = () => {
    let currentFormIsValid = true;

    //Name
    if (!this.name) {
      currentFormIsValid = false;
      this.setErrors('name', 'Необхідно заповнити');
    }

    if (this.name) {
      if (!this.name.match(/^[a-zA-Z]+$/)) {
        currentFormIsValid = false;
        this.setErrors('name', 'Має містити тільки букви');
      }
    }

     //Phone
     if (!this.phone) {
      currentFormIsValid = false;
      this.setErrors('phone', 'Необхідно заповнити');
    }

    if (this.phone) {
      if (!this.phone.match(/^[\+]?[(]?[0-9]{3}[)]?[-\s\.]?[0-9]{3}[-\s\.]?[0-9]{4,6}$/im)) {
        currentFormIsValid = false;
        this.setErrors('phone', 'Має бути номер телефону');
      }
    }

    //Email
    if (!this.email) {
      currentFormIsValid = false;
      this.setErrors('email', 'Необхідно заповнити');
    }

    if (this.email) {
      if (!this.email.match(/[a-zA-Z0-9]+[\.]?([a-zA-Z0-9]+)?[\@][a-z]{3,9}[\.][a-z]{2,5}/g)) {
        currentFormIsValid = false;
        this.setErrors('email', 'Має бути у форматі email');
      }
    }

     //Description
     if (!this.description) {
      currentFormIsValid = false;
      this.setErrors('description', 'Необхідно заповнити');
    }

    return currentFormIsValid;
  }
}

export default new ContactsBlockService();
