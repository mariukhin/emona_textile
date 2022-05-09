// node modules
import React, { useRef, useEffect, RefObject } from 'react';
// components
import BlockInfoComponent from 'components/BlockInfoComponent';
import ContactsBlock from 'components/Footer/components/ContactsBlock';
import { Wrapper } from '@googlemaps/react-wrapper';
import { Box, FormControl, InputLabel } from '@mui/material';
// styles
import {
  ContactsAndFormBlockWrapper,
  BlockContainer,
  StyledMapComponent,
  InfoContainer,
  StyledPaper,
  PaperWrapper,
  FormHeader,
  StyledOutlinedInput,
  StyledButton,
  StyledButtonText,
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

  return <StyledMapComponent ref={ref as unknown as RefObject<HTMLDivElement>} id="map" />;
};

const ContactsAndFormBlock = () => {
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

        <StyledPaper>
          <PaperWrapper>
            <FormHeader>Заповніть форму</FormHeader>

            <Box
              component="form"
              autoComplete="off"
            >
              <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: '20px' }}>
                <FormControl style={{ width: '48%' }}>
                  <InputLabel htmlFor="name">Ім'я</InputLabel>
                  <StyledOutlinedInput
                    id="name"
                    margin="dense"
                    type="text"
                    label="Ім'я"
                  />
                </FormControl>
                <FormControl style={{ width: '48%' }}>
                  <InputLabel htmlFor="phone">Телефон</InputLabel>
                  <StyledOutlinedInput
                    id="phone"
                    margin="dense"
                    type="tel"
                    label="Телефон"
                  />
                </FormControl>
              </div>
              <FormControl style={{ marginTop: '12px', width: '100%' }}>
                <InputLabel htmlFor="email">Електронна пошта</InputLabel>
                <StyledOutlinedInput
                  id="email"
                  label="Електронна пошта"
                  fullWidth
                  margin="dense"
                  type="email"
                />
              </FormControl>
              <FormControl style={{ margin: '12px 0', width: '100%' }}>
                <InputLabel htmlFor="name">Опис замовлення</InputLabel>
                <StyledOutlinedInput
                  id="description"
                  multiline
                  rows={6}
                  label="Опис замовлення"
                  fullWidth
                  margin="dense"
                  type="text"
                />
              </FormControl>

              <StyledButton
                color="success"
                variant="contained"
                size="small"
              >
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

export default ContactsAndFormBlock;
