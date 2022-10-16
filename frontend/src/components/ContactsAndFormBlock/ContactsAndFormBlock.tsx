// node modules
import * as R from 'ramda';
import React, { useRef, useEffect, RefObject } from 'react';
import { observer } from 'mobx-react';
// components
import BlockInfoComponent from 'components/BlockInfoComponent';
import ContactsBlock from 'components/Footer/components/ContactsBlock';
import { Wrapper } from '@googlemaps/react-wrapper';
import { Box, InputLabel, OutlinedInput, FormHelperText } from '@mui/material';
import { useStore } from 'modules/Stores';
// styles
import {
  ContactsAndFormBlockWrapper,
  BlockContainer,
  StyledMapComponent,
  InfoContainer,
  StyledPaper,
  PaperWrapper,
  FormHeader,
  StyledButton,
  StyledButtonText,
  StyledFormControl,
} from './styles';
import { colors } from 'utils/color';

const MapComponent = () => {
  const ref = useRef();

  useEffect(() => {
    if (ref.current) {
      new window.google.maps.Map(ref.current, {
        center: { lat: 50.454906, lng: 30.488335 },
        zoom: 17,
      });
    }
  });

  return (
    <StyledMapComponent
      ref={ref as unknown as RefObject<HTMLDivElement>}
      id="map"
    />
  );
};

const ContactsAndFormBlock = () => {
  const {
    name,
    phone,
    email,
    description,
    setName,
    setPhone,
    setEmail,
    setDescription,
    errors,
    clearErrorByKey,
    handleValidateForm,
    sendEmail,
  } = useStore('ContactsAndFormBlockStore');

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { id, value } = event.target;

    switch (id) {
      case 'name':
        return setName(value);
      case 'phone':
        return setPhone(value);
      case 'email':
        return setEmail(value);
      case 'description':
        return setDescription(value);
      default:
        break;
    }
  };

  const handleFocus = (event: React.FocusEvent<HTMLInputElement>) => {
    clearErrorByKey(event.target.id);
  };

  const handleSubmit = (event: React.FormEvent<HTMLButtonElement>) => {
    event.preventDefault();
    const formIsValid = handleValidateForm();

    if (formIsValid) {
      sendEmail();
    }
  }

  return (
    <ContactsAndFormBlockWrapper>
      <BlockInfoComponent title="Контакти" subtitle="Зв'язатися з нами" />

      <BlockContainer>
        <InfoContainer>
          <ContactsBlock />

          <Wrapper apiKey="AIzaSyBCpmLQv2zdquHe3Pk2Jh_qNpscEkkKhWE">
            <MapComponent />
          </Wrapper>
        </InfoContainer>

        <StyledPaper sx={{ height: !R.isEmpty(errors) ? '550px' : '484px' }}>
          <PaperWrapper>
            <FormHeader>Заповніть форму</FormHeader>

            <Box component="form" autoComplete="off">
              <div
                style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  marginTop: '20px',
                }}
              >
                <StyledFormControl style={{ width: '48%' }}>
                  <InputLabel htmlFor="name">Ім'я</InputLabel>
                  <OutlinedInput
                    id="name"
                    value={name}
                    onChange={handleChange}
                    onFocus={handleFocus}
                    margin="dense"
                    type="text"
                    label="Ім'я"
                    error={!!errors.name}
                  />
                  {errors.name && <FormHelperText error>{errors.name}</FormHelperText>}
                </StyledFormControl>
                <StyledFormControl style={{ width: '48%' }}>
                  <InputLabel htmlFor="phone">Телефон</InputLabel>
                  <OutlinedInput
                    id="phone"
                    value={phone}
                    onChange={handleChange}
                    onFocus={handleFocus}
                    margin="dense"
                    type="tel"
                    label="Телефон"
                    error={!!errors.phone}
                  />
                  {errors.phone && <FormHelperText error>{errors.phone}</FormHelperText>}
                </StyledFormControl>
              </div>
              <StyledFormControl style={{ marginTop: '12px', width: '100%' }}>
                <InputLabel htmlFor="email">Електронна пошта</InputLabel>
                <OutlinedInput
                  id="email"
                  value={email}
                  onChange={handleChange}
                  onFocus={handleFocus}
                  label="Електронна пошта"
                  fullWidth
                  margin="dense"
                  type="email"
                  error={!!errors.email}
                />
                {errors.email && <FormHelperText error>{errors.email}</FormHelperText>}
              </StyledFormControl>
              <StyledFormControl style={{ margin: '12px 0', width: '100%' }}>
                <InputLabel htmlFor="name">Опис замовлення</InputLabel>
                <OutlinedInput
                  id="description"
                  value={description}
                  onChange={handleChange}
                  onFocus={handleFocus}
                  multiline
                  rows={6}
                  label="Опис замовлення"
                  fullWidth
                  margin="dense"
                  type="text"
                  error={!!errors.description}
                />
                {errors.description && <FormHelperText error>{errors.description}</FormHelperText>}
              </StyledFormControl>

              <StyledButton color="success" variant="contained" size="small" onClick={handleSubmit}>
                <StyledButtonText variant="button" color={colors.text.white}>
                  Надіслати
                </StyledButtonText>
              </StyledButton>
            </Box>
          </PaperWrapper>
        </StyledPaper>
      </BlockContainer>
    </ContactsAndFormBlockWrapper>
  );
};

export default observer(ContactsAndFormBlock);
