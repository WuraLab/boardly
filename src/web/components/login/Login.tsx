/* eslint-disable @typescript-eslint/explicit-function-return-type */
/* eslint-disable @typescript-eslint/explicit-module-boundary-types */
import Box from '@material-ui/core/Box';
import Grid from '@material-ui/core/Grid';
import { Container, Button, Typography } from '@material-ui/core';
import '@fontsource/poppins';
import { FormInput, FormInputPassword } from './FormInput';

const Login = () => {
    return (
        <Box padding={0} margin={0}>
            <Grid container spacing={0}>
                <Grid item xs={6}>
                    <Box style={{ margin: '20px 0 20px 100px' }}>
                        <img src="/logo.png" width="100" height="100" alt="" />
                    </Box>
                    <Container
                        style={{
                            height: '100vh',
                            backgroundColor: 'white',
                            display: 'flex',
                            margin: '50px 0 0 0'
                        }}>
                        <Grid
                            item
                            container
                            direction="column"
                            style={{ width: '60%', margin: '0 auto' }}>
                            <Typography
                                style={{
                                    color: '#000000',
                                    fontWeight: 'bold',
                                    fontSize: '30px',
                                    letterSpacing: '4px'
                                }}>
                                Sign in to Boardly
                            </Typography>

                            <br />

                            <FormInput name="Company's Email" placeholder="contact@wuralab.com">
                                {' '}
                            </FormInput>
                            <FormInputPassword name="Password" placeholder="Password">
                                {' '}
                            </FormInputPassword>
                            <Typography color="primary">
                                <a href="/create-admin">Forget Password</a>
                            </Typography>

                            <br />
                            <Button
                                variant="contained"
                                style={{
                                    backgroundColor: '#22ad80',
                                    color: '#ffffff'
                                }}>
                                Sign in
                            </Button>
                            <Typography color="primary">
                                New here?
                                <a href="/create-admin-view">
                                    <span> Create an account </span>
                                </a>
                            </Typography>
                        </Grid>
                    </Container>
                </Grid>
                <Grid item xs={6}>
                    <Box>
                        <Container
                            style={{
                                height: '100vh',
                                backgroundColor: '#753ff6',
                                clipPath: 'polygon(8% 0, 100% 0, 100% 100%, 0% 100%)'
                            }}>
                            <div
                                style={{
                                    margin: '0 50px 0 50px',
                                    position: 'absolute',
                                    left: '50%',
                                    top: '50%',
                                    transform: 'translate(10%, -50%)'
                                }}>
                                <img src="/admin-signup.png" alt="" />
                            </div>
                        </Container>
                    </Box>
                </Grid>
            </Grid>
        </Box>
    );
};

export default Login;
