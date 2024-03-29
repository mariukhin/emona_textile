// node modules
import * as R from 'ramda';
import React, { useRef, useEffect, RefObject } from 'react';
import { observer } from 'mobx-react';
import { IMaskInput } from 'react-imask';
// components
import BlockInfoComponent from 'components/BlockInfoComponent';
import ContactsBlock from 'components/Footer/components/ContactsBlock';
import { Wrapper } from '@googlemaps/react-wrapper';
import { Box, InputLabel, FormHelperText } from '@mui/material';
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
  FormBlock,
  StyledButton,
  StyledButtonText,
  StyledFormControl,
  StyledOutlinedInput
} from './styles';
import { colors } from 'utils/color';

export const addSingleMarkers = ({
  locations,
  map,
}: {
  locations: ReadonlyArray<google.maps.LatLngLiteral>;
  map: google.maps.Map | null | undefined;
}) =>
  locations.map(
    position =>
      new google.maps.Marker({
        position,
        map,
      }),
  );

const MapComponent = () => {
  const ref = useRef();

  useEffect(() => {
    if (ref.current) {
      const map = new window.google.maps.Map(ref.current, {
        center: { lat: 50.4549785, lng: 30.4884064 },
        zoom: 17,
      });

      new google.maps.Marker({
        position: { lat: 50.45499, lng: 30.48842 },
        map,
      })
    }
  });

  return (
    <StyledMapComponent
      ref={ref as unknown as RefObject<HTMLDivElement>}
      id="map"
    />
  );
};

const TextMaskCustom = React.forwardRef<HTMLInputElement, {
  onChange: (event: { target: { id: string; value: string } }) => void;
  name: string;
}>(
  function TextMaskCustom(props, ref) {
    const { onChange, ...other } = props;
    return (
      <IMaskInput
        {...other}
        mask="+38\0 (00) - 000 - 00 - 00"
        definitions={{
          '#': /[1-9]/,
        }}
        lazy={false}
        inputRef={ref}
        onAccept={(value: any) => onChange({ target: { id: 'phone', value } })}
        overwrite
      />
    );
  },
);

const ContactsAndFormBlock = ({isCatalogItemPage = false}) => {
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

  const blockTitle = isCatalogItemPage ? 'Замовити' : 'Контакти';

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
      <BlockInfoComponent title={ blockTitle } subtitle="Зв'язатися з нами" />

      <BlockContainer>
        <InfoContainer>
          <ContactsBlock />

          <Wrapper apiKey="AIzaSyDXL3kpn2D415JAefS9_tYo-Jvz6Sb0ZRQ">
            <MapComponent />
          </Wrapper>
        </InfoContainer>

        <StyledPaper isErrors={ !R.isEmpty(errors) } id="contact-form-anchor">
          <PaperWrapper>
            <FormHeader>Заповніть форму</FormHeader>

            <Box component="form" autoComplete="off">
              <FormBlock>
                <StyledFormControl width={'48%'}>
                  <InputLabel htmlFor="name">Ім'я</InputLabel>
                  <StyledOutlinedInput
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
                <StyledFormControl width={'48%'}>
                  <StyledOutlinedInput
                    id="phone"
                    value={phone}
                    onChange={handleChange}
                    onFocus={handleFocus}
                    inputComponent={TextMaskCustom as any}
                    type="tel"
                    error={!!errors.phone}
                  />
                  {errors.phone && <FormHelperText error>{errors.phone}</FormHelperText>}
                </StyledFormControl>
              </FormBlock>
              <StyledFormControl width={'100%'} marginCustom={'12px 0 0'}>
                <InputLabel htmlFor="email">Електронна пошта</InputLabel>
                <StyledOutlinedInput
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
              <StyledFormControl width={'100%'} marginCustom={'12px 0'}>
                <InputLabel htmlFor="name">Опис замовлення</InputLabel>
                <StyledOutlinedInput
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
