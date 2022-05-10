// store
import ContactsBlockStore from './store';

class ContactsBlockService extends ContactsBlockStore {
  handleValidateForm = () => {
    let currentFormIsValid = true;

    console.log('Here');

    //Name
    if (!this.name) {
      currentFormIsValid = false;
      console.log('Here1');
      this.setErrors('name', 'Необхідно заповнити')
    }

    if (this.name) {
      if (!this.name.match(/^[a-zA-Z]+$/)) {
        currentFormIsValid = false;
        this.setErrors('name', 'Має містити тільки букви')
      }
    }

    //Email
    if (!this.email) {
      currentFormIsValid = false;
      console.log('Here2');
      this.setErrors('email', 'Необхідно заповнити')
    }

    if (this.email) {
      if (!this.email.match(/[a-zA-Z0-9]+[\.]?([a-zA-Z0-9]+)?[\@][a-z]{3,9}[\.][a-z]{2,5}/g)) {
        currentFormIsValid = false;
        this.setErrors('email', 'Має бути у форматі email')
      }
    }

    console.log('currentFormIsValid', currentFormIsValid);

    this.setFormIsValid(currentFormIsValid);
  }
}

export default new ContactsBlockService();
