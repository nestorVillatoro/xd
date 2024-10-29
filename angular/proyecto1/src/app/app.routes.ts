import { Routes } from '@angular/router';
import { LoginComponent } from './login/login.component';
import { VisualizadorDiscosComponent } from './explorador-archivos/explorador-archivos.component';
import { ConsolaLogeadaComponent } from './consola-logeada/consola-logeada.component';

export const routes: Routes = [
    
    { path: 'login', component: LoginComponent },
     {path: 'explorador-archivos', component: VisualizadorDiscosComponent},
    { path: 'consola-logeada', component: ConsolaLogeadaComponent},
    { path: '**', redirectTo: '' }
];