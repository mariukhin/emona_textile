// node modules
import React from 'react';
// styles
import {
  PagePhotoBlockContainer,
  Heading,
  InfoBlock,
  StyledButton,
  StyledButtonText,
} from './styles';

type PagePhotoBlockProps = {
  heading: string;
  btnText: string;
  imageUrl: string;
};

const PagePhotoBlock: React.FC<PagePhotoBlockProps> = ({
  heading,
  btnText,
  imageUrl,
}) => (
  <PagePhotoBlockContainer imageUrl={ imageUrl }>
    <InfoBlock>
      <Heading>{heading}</Heading>
      <StyledButton color="warning" size="large" variant="contained">
        <StyledButtonText>{btnText}</StyledButtonText>
      </StyledButton>
    </InfoBlock>
  </PagePhotoBlockContainer>
);

export default PagePhotoBlock;
